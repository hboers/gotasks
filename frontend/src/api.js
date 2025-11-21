const API_BASE = 'http://localhost:8080/api'

export async function fetchTodos() {
  const res = await fetch(`${API_BASE}/todos`)
  if (!res.ok) throw new Error('Failed to load todos')
  return res.json()
}

export async function createTodo(title) {
  const res = await fetch(`${API_BASE}/todos`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title })
  })
  if (!res.ok) throw new Error('Failed to create todo')
  return res.json()
}

export async function updateTodo(id, patch) {
  const res = await fetch(`${API_BASE}/todos/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(patch)
  })
  if (!res.ok) throw new Error('Failed to update todo')
  return res.json()
}

export async function deleteTodo(id) {
  const res = await fetch(`${API_BASE}/todos/${id}`, {
    method: 'DELETE'
  })
  if (!res.ok && res.status !== 204) {
    throw new Error('Failed to delete todo')
  }
}
