const API_BASE = "http://localhost:8080";

const state = {
  mode: "login",
  accessToken: localStorage.getItem("access_token") || "",
  me: null,
  chats: [],
  activeChat: null,
  ws: null,
  reconnectTimer: null,
  typingStopTimer: null,
  lastTypingSent: false
};

const $ = id => document.getElementById(id);
const els = {
  toast:$("toast"), authView:$("authView"), appView:$("appView"),
  loginTab:$("loginTab"), registerTab:$("registerTab"), authForm:$("authForm"),
  usernameInput:$("usernameInput"), emailWrap:$("emailWrap"), emailInput:$("emailInput"),
  passwordInput:$("passwordInput"), authSubmit:$("authSubmit"), meLabel:$("meLabel"),
  logoutBtn:$("logoutBtn"), userSearchInput:$("userSearchInput"), userSearchResults:$("userSearchResults"),
  chatList:$("chatList"), emptyState:$("emptyState"), chatView:$("chatView"),
  chatTitle:$("chatTitle"), chatPresence:$("chatPresence"), messages:$("messages"),
  typingIndicator:$("typingIndicator"), messageForm:$("messageForm"), messageInput:$("messageInput"),
  reloadBtn:$("reloadBtn")
};

function toast(msg){
  els.toast.textContent = msg;
  els.toast.classList.remove("hidden");
  clearTimeout(toast.t);
  toast.t = setTimeout(()=>els.toast.classList.add("hidden"), 3500);
}

function setMode(mode){
  state.mode = mode;
  const login = mode === "login";
  els.loginTab.classList.toggle("active", login);
  els.registerTab.classList.toggle("active", !login);
  els.emailWrap.classList.toggle("hidden", login);
  els.emailInput.required = !login;
  els.authSubmit.textContent = login ? "Войти" : "Зарегистрироваться";
}

function extractToken(p){
  return p?.access_token || p?.accessToken || p?.token || p?.tokens?.access_token || p?.data?.access_token || "";
}

async function request(path, options={}){
  const headers = new Headers(options.headers || {});
  if (options.body != null && !headers.has("Content-Type")) headers.set("Content-Type","application/json");
  if (state.accessToken) headers.set("Authorization", `Bearer ${state.accessToken}`);

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
    credentials:"include"
  });

  if (res.status === 204) return null;

  const text = await res.text();
  let data = null;
  try { data = text ? JSON.parse(text) : null; } catch { data = text; }

  if (!res.ok) throw new Error(data?.error || data?.message || data || `HTTP ${res.status}`);
  return data;
}

async function login(username,password){
  const data = await request("/auth/login", {
    method:"POST",
    body:JSON.stringify({username,password})
  });
  const token = extractToken(data);
  if (!token) throw new Error("Backend не вернул access_token");
  state.accessToken = token;
  localStorage.setItem("access_token", token);
}

async function register(username,email,password){
  const data = await request("/auth/register", {
    method:"POST",
    body:JSON.stringify({username,email,password})
  });
  const token = extractToken(data);
  if (token) {
    state.accessToken = token;
    localStorage.setItem("access_token", token);
  } else {
    await login(username,password);
  }
}

async function loadMe(){
  state.me = await request("/api/me");
  const u = state.me?.username || state.me?.user?.username || `user #${state.me?.id ?? "?"}`;
  els.meLabel.textContent = u;
}

async function loadChats(){
  const data = await request("/api/chats");
  state.chats = Array.isArray(data) ? data : (data?.chats || []);
  renderChats();
  if (state.activeChat) {
    const fresh = state.chats.find(c => Number(c.id) === Number(state.activeChat.id));
    if (fresh) {
      state.activeChat = fresh;
      renderHeader();
    }
  }
}

function formatLastSeen(v){
  if (!v) return "не в сети";
  const d = new Date(v);
  return Number.isNaN(d.getTime()) ? "не в сети" : `был(а) в сети ${d.toLocaleString("ru-RU")}`;
}

function renderChats(){
  els.chatList.innerHTML = "";
  if (!state.chats.length) {
    els.chatList.innerHTML = '<div class="muted small-text" style="padding:12px">Чатов пока нет</div>';
    return;
  }

  for (const chat of state.chats) {
    const item = document.createElement("div");
    item.className = "chat-item";
    if (state.activeChat && Number(state.activeChat.id) === Number(chat.id)) item.classList.add("active");

    const meta = document.createElement("div");
    const name = document.createElement("div");
    name.className = "chat-name";
    name.textContent = chat.title || `Chat #${chat.id}`;

    const status = document.createElement("div");
    status.className = "status";
    status.textContent = chat.is_online ? "online" : formatLastSeen(chat.last_seen_at);

    meta.append(name,status);
    item.append(meta);

    if (chat.is_online) {
      const dot = document.createElement("div");
      dot.className = "online-dot";
      item.append(dot);
    }

    item.onclick = ()=>openChat(chat);
    els.chatList.append(item);
  }
}

async function openChat(chat){
  state.activeChat = chat;
  renderChats();
  renderHeader();
  els.emptyState.classList.add("hidden");
  els.chatView.classList.remove("hidden");
  await loadMessages(chat.id);
  try { await request(`/api/chats/${chat.id}/messages/read`, {method:"PATCH"}); } catch {}
}

function renderHeader(){
  if (!state.activeChat) return;
  els.chatTitle.textContent = state.activeChat.title || `Chat #${state.activeChat.id}`;
  els.chatPresence.textContent = state.activeChat.is_online ? "online" : formatLastSeen(state.activeChat.last_seen_at);
}

function myID(){
  return Number(state.me?.user_id ?? 0);
}

async function loadMessages(chatID){
  const data = await request(`/api/chats/${chatID}/messages`);
  const list = Array.isArray(data) ? data : (data?.messages || []);
  els.messages.innerHTML = "";
  for (const m of list) appendMessage(m);
  els.messages.scrollTop = els.messages.scrollHeight;
}

function appendMessage(m){
  const box = document.createElement("div");
  const mine = Number(m.sender_id ?? m.user_id ?? 0) === myID();
  box.className = `message${mine ? " mine" : ""}`;

  const text = document.createElement("div");
  text.textContent = m.text ?? "";

  const info = document.createElement("div");
  info.className = "message-info";
  const d = m.sent_at ? new Date(m.sent_at) : null;
  const time = d && !Number.isNaN(d.getTime()) ? d.toLocaleTimeString("ru-RU",{hour:"2-digit",minute:"2-digit"}) : "";
  info.textContent = [time,m.status].filter(Boolean).join(" · ");

  box.append(text,info);
  els.messages.append(box);
}

async function sendMessage(){
  if (!state.activeChat) return;
  const text = els.messageInput.value.trim();
  if (!text) return;
  const data = await request(`/api/chats/${state.activeChat.id}/messages`, {
    method:"POST",
    body:JSON.stringify({text})
  });
  const m = data?.message;
  if (m) appendMessage(m);
  els.messageInput.value = "";
  els.messages.scrollTop = els.messages.scrollHeight;
  sendTyping(false);
}

let searchTimer;
async function searchUsers(){
  const q = els.userSearchInput.value.trim();
  els.userSearchResults.innerHTML = "";
  if (!q) return;

  try {
    const e = encodeURIComponent(q);
    const data = await request(`/api/users/search?username=${e}`);
    const users = Array.isArray(data) ? data : (data?.users || []);

    for (const user of users) {
      if (Number(user.id) === myID()) continue;

      const row = document.createElement("div");
      row.className = "search-result";

      const name = document.createElement("span");
      name.textContent = user.username || `user #${user.id}`;

      const btn = document.createElement("button");
      btn.className = "ghost";
      btn.type = "button";
      btn.textContent = "Чат";
      btn.onclick = ()=>createChat(user);

      row.append(name,btn);
      els.userSearchResults.append(row);
    }
  } catch(err) {
    toast(`Поиск: ${err.message}`);
  }
}

async function createChat(user){
  await request("/api/chats", {
    method:"POST",
    body:JSON.stringify({ username:user.username })
  });
  els.userSearchInput.value = "";
  els.userSearchResults.innerHTML = "";
  await loadChats();
}

function wsURL(){
  const u = new URL(API_BASE);
  u.protocol = u.protocol === "https:" ? "wss:" : "ws:";
  u.pathname = "/api/ws";
  u.search = `?token=${encodeURIComponent(state.accessToken)}`;
  return u.toString();
}

function connectWS(){
  if (!state.accessToken) return;
  if (state.ws) try { state.ws.close(); } catch {}

  const socket = new WebSocket(wsURL());
  state.ws = socket;

  socket.onmessage = async e => {
    let event;
    try { event = JSON.parse(e.data); } catch { return; }
    const d = event?.data || {};

    if (event.type === "new_message") {
      if (state.activeChat && Number(state.activeChat.id) === Number(d.chat_id)) {
        appendMessage(d);
        els.messages.scrollTop = els.messages.scrollHeight;
        try { await request(`/api/chats/${d.chat_id}/messages/read`, {method:"PATCH"}); } catch {}
      }
    }

    if (event.type === "messages_read" && state.activeChat && Number(state.activeChat.id) === Number(d.chat_id)) {
      loadMessages(state.activeChat.id).catch(()=>{});
    }

    if (event.type === "typing" && state.activeChat && Number(state.activeChat.id) === Number(d.chat_id)) {
      els.typingIndicator.classList.toggle("hidden", !d.is_typing);
    }

    if (event.type === "user_online" || event.type === "user_offline") {
      loadChats().catch(()=>{});
    }
      if (event.type === "chat_created") {
          await loadChats();
      }
  };

  socket.onclose = () => {
    if (state.ws === socket) state.ws = null;
    if (state.accessToken) {
      clearTimeout(state.reconnectTimer);
      state.reconnectTimer = setTimeout(connectWS,1500);
    }
  };
}

function sendTyping(isTyping){
  if (!state.activeChat || !state.ws || state.ws.readyState !== WebSocket.OPEN) return;
  if (state.lastTypingSent === isTyping) return;
  state.lastTypingSent = isTyping;
  state.ws.send(JSON.stringify({
    type:"typing",
    data:{chat_id:Number(state.activeChat.id),is_typing:isTyping}
  }));
}

async function enterApp(){
  await loadMe();
  await loadChats();
  els.authView.classList.add("hidden");
  els.appView.classList.remove("hidden");
  connectWS();
}

async function logout(){
  try { await request("/auth/logout",{method:"POST"}); } catch {}
  state.accessToken = "";
  localStorage.removeItem("access_token");
  if (state.ws) try { state.ws.close(); } catch {}
  els.appView.classList.add("hidden");
  els.authView.classList.remove("hidden");
}

els.loginTab.onclick = ()=>setMode("login");
els.registerTab.onclick = ()=>setMode("register");

els.authForm.onsubmit = async e => {
  e.preventDefault();
  const username = els.usernameInput.value.trim();
  const email = els.emailInput.value.trim();
  const password = els.passwordInput.value;

  try {
    if (state.mode === "login") await login(username,password);
    else await register(username,email,password);
    await enterApp();
  } catch(err) {
    toast(err.message);
  }
};

els.logoutBtn.onclick = logout;
els.userSearchInput.oninput = () => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(searchUsers,250);
};
els.messageForm.onsubmit = async e => {
  e.preventDefault();
  try { await sendMessage(); } catch(err) { toast(err.message); }
};
els.messageInput.oninput = () => {
  sendTyping(true);
  clearTimeout(state.typingStopTimer);
  state.typingStopTimer = setTimeout(()=>sendTyping(false),900);
};
els.messageInput.onblur = ()=>sendTyping(false);
els.reloadBtn.onclick = async () => {
  try {
    await loadChats();
    if (state.activeChat) await loadMessages(state.activeChat.id);
  } catch(err) {
    toast(err.message);
  }
};

setMode("login");
if (state.accessToken) enterApp().catch(()=>logout());
