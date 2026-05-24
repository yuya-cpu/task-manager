const API_BASE = window.location.origin;

const authSection = document.getElementById("auth-section");
const appSection = document.getElementById("app-section");
const authMessage = document.getElementById("auth-message");
const taskMessage = document.getElementById("task-message");
const taskList = document.getElementById("task-list");
const userLabel = document.getElementById("user-label");

function getToken() {
  return localStorage.getItem("token");
}

function setToken(token) {
  localStorage.setItem("token", token);
}

function setUserEmail(email) {
  localStorage.setItem("email", email);
}

function clearToken() {
  localStorage.removeItem("token");
  localStorage.removeItem("email");
}

function showMessage(el, text, ok = false) {
  el.textContent = text;
  el.classList.toggle("ok", ok);
}

async function api(path, options = {}) {
  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };
  const token = getToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  if (res.status === 204) {
    return null;
  }

  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return body;
}

function showApp(email) {
  authSection.classList.add("hidden");
  appSection.classList.remove("hidden");
  userLabel.textContent = `ログイン中: ${email}`;
  loadTasks();
}

function showAuth() {
  appSection.classList.add("hidden");
  authSection.classList.remove("hidden");
}

async function login() {
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;
  try {
    const data = await api("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    setToken(data.token);
    setUserEmail(email);
    showMessage(authMessage, "ログインしました", true);
    showApp(email);
  } catch (err) {
    showMessage(authMessage, err.message);
  }
}

async function signup() {
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;
  try {
    await api("/auth/signup", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    showMessage(authMessage, "登録しました。ログインしてください。", true);
  } catch (err) {
    showMessage(authMessage, err.message);
  }
}

function buildListQuery() {
  const params = new URLSearchParams();
  const status = document.getElementById("filter-status")?.value;
  const priority = document.getElementById("filter-priority")?.value;
  const sort = document.getElementById("filter-sort")?.value;
  if (status) params.set("status", status);
  if (priority) params.set("priority", priority);
  if (sort) params.set("sort", sort);
  const qs = params.toString();
  return qs ? `/assignments?${qs}` : "/assignments";
}

async function loadTasks() {
  taskList.innerHTML = "";
  const listMeta = document.getElementById("list-meta");
  if (listMeta) listMeta.textContent = "";
  try {
    const data = await api(buildListQuery());
    if (listMeta && data.meta) {
      listMeta.textContent = `全 ${data.meta.total} 件（${data.meta.page} ページ目 / ${data.meta.limit} 件表示）`;
    }
    if (!data.data.length) {
      taskList.innerHTML = "<li>タスクがありません</li>";
      return;
    }
    data.data.forEach((task) => {
      const li = document.createElement("li");
      li.className = "task-item";
      const due = task.due_date
        ? new Date(task.due_date).toLocaleDateString("ja-JP")
        : "なし";
      li.innerHTML = `
        <h3>${escapeHtml(task.title)}</h3>
        <p>${escapeHtml(task.description || "")}</p>
        <div class="task-meta">
          <span class="badge">${task.priority}</span>
          <span class="badge">${task.status}</span>
          <span>期限: ${due}</span>
        </div>
        <div class="task-actions">
          <button data-action="done" data-id="${task.id}">完了</button>
          <button data-action="delete" data-id="${task.id}" class="secondary">削除</button>
        </div>
      `;
      taskList.appendChild(li);
    });
  } catch (err) {
    showMessage(taskMessage, err.message);
    if (err.message.includes("401") || err.message.toLowerCase().includes("unauthorized")) {
      clearToken();
      showAuth();
    }
  }
}

function escapeHtml(text) {
  return text
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

async function createTask(event) {
  event.preventDefault();
  const payload = {
    title: document.getElementById("title").value,
    description: document.getElementById("description").value,
    priority: document.getElementById("priority").value,
    status: document.getElementById("status").value,
  };
  const dueDate = document.getElementById("due-date").value;
  if (dueDate) {
    payload.due_date = dueDate;
  }

  try {
    await api("/assignments", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    document.getElementById("task-form").reset();
    showMessage(taskMessage, "タスクを追加しました", true);
    loadTasks();
  } catch (err) {
    showMessage(taskMessage, err.message);
  }
}

async function markDone(id) {
  try {
    await api(`/assignments/${id}`, {
      method: "PUT",
      body: JSON.stringify({ status: "done" }),
    });
    loadTasks();
  } catch (err) {
    showMessage(taskMessage, err.message);
  }
}

async function deleteTask(id) {
  try {
    await api(`/assignments/${id}`, { method: "DELETE" });
    loadTasks();
  } catch (err) {
    showMessage(taskMessage, err.message);
  }
}

document.getElementById("login-btn").addEventListener("click", login);
document.getElementById("signup-btn").addEventListener("click", signup);
document.getElementById("logout-btn").addEventListener("click", () => {
  clearToken();
  showAuth();
});
document.getElementById("task-form").addEventListener("submit", createTask);
document.getElementById("filter-btn").addEventListener("click", loadTasks);

taskList.addEventListener("click", (event) => {
  const button = event.target.closest("button");
  if (!button) return;
  const id = button.dataset.id;
  if (button.dataset.action === "done") {
    markDone(id);
  }
  if (button.dataset.action === "delete") {
    deleteTask(id);
  }
});

if (getToken()) {
  showApp(localStorage.getItem("email") || "ユーザー");
} else {
  showAuth();
}
