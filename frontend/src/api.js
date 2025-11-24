const API_BASE = 'http://localhost:8080/api'

export async function register(email, name, password) {
  const res = await fetch(`${API_BASE}/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, name, password }),
  })

  if (!res.ok) {
    throw new Error(await res.text())
  }

  return res.json() // { id, email, name }
}

export async function login(email, password) {
  const res = await fetch(`${API_BASE}/login`, {
    method: "POST",
    credentials: "include",              // 🔥 Cookie senden & akzeptieren
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) {
    const text = await res.text();
    const err = new Error(text || "Login failed");
    err.status = res.status;
    throw err;
  }

  return res.json();
}

export async function fetchMe() {
  const res = await fetch(`${API_BASE}/me`, {
    credentials: "include"  
  })

  if (!res.ok) {
    const err = new Error("Unauthorized")
    err.status = res.status
    throw err
  }
  return res.json() // { id, email, name }
}

export async function fetchTodos() {
  const res = await fetch(`${API_BASE}/todos`,{
    credentials: "include"  
  }) 
  if (!res.ok) {
    const err = new Error('Failed to load todos')
    err.status = res.status
    throw  err
  }
  return res.json()
}

export async function createTodo(title) {
  const res = await fetch(`${API_BASE}/todos`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title }),
    credentials: "include"  
  })
  if (!res.ok) throw new Error('Failed to create todo')
  return res.json()
}

export async function updateTodo(id, done) {
  const res = await fetch(`${API_BASE}/todos/${id}`, {
    method: "PUT",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ done }),
  });

  if (!res.ok) {
    const t = await res.text();
    const err = new Error(t);
    err.status = res.status;
    throw err;
  }

  return res.json();
}

export async function deleteTodo(id) {
  const res = await fetch(`${API_BASE}/todos/${id}`, {
    method: 'DELETE',
    credentials: "include"  
  })
  if (!res.ok && res.status !== 204) {
    throw new Error('Failed to delete todo')
  }
}
