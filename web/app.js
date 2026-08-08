const API_BASE = "/api/clients";

const form = document.getElementById("client-form");
const formTitle = document.getElementById("form-title");
const btnSubmit = document.getElementById("btn-submit");
const btnCancel = document.getElementById("btn-cancel");
const clientIdInput = document.getElementById("client-id");
const nombreInput = document.getElementById("nombre");
const emailInput = document.getElementById("email");
const telefonoInput = document.getElementById("telefono");
const direccionInput = document.getElementById("direccion");

const loadingEl = document.getElementById("loading");
const errorEl = document.getElementById("error");
const emptyEl = document.getElementById("empty");
const tableEl = document.getElementById("clients-table");
const tbody = document.getElementById("clients-body");
const btnRefresh = document.getElementById("btn-refresh");

let editingId = null;

async function fetchClients() {
  showLoading(true);
  hideError();
  try {
    const res = await fetch(API_BASE);
    if (!res.ok) throw new Error("Error al cargar clientes");
    const data = await res.json();
    renderClients(data || []);
  } catch (err) {
    showError(err.message);
  } finally {
    showLoading(false);
  }
}

function renderClients(clients) {
  tbody.innerHTML = "";
  if (!clients.length) {
    tableEl.classList.add("hidden");
    emptyEl.classList.remove("hidden");
    return;
  }
  emptyEl.classList.add("hidden");
  tableEl.classList.remove("hidden");

  clients.forEach((c) => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${escapeHtml(c.nombre)}</td>
      <td>${escapeHtml(c.email)}</td>
      <td>${escapeHtml(c.telefono || "-")}</td>
      <td>${escapeHtml(c.direccion || "-")}</td>
      <td class="actions">
        <button class="btn btn-edit" data-id="${c.id}">Editar</button>
        <button class="btn btn-danger" data-id="${c.id}">Eliminar</button>
      </td>
    `;
    tbody.appendChild(tr);
  });

  // Event listeners
  tbody.querySelectorAll(".btn-edit").forEach((btn) => {
    btn.addEventListener("click", () => startEdit(btn.dataset.id, clients));
  });
  tbody.querySelectorAll(".btn-danger").forEach((btn) => {
    btn.addEventListener("click", () => deleteClient(btn.dataset.id));
  });
}

function startEdit(id, clients) {
  const client = clients.find((c) => c.id === id);
  if (!client) return;

  editingId = id;
  clientIdInput.value = id;
  nombreInput.value = client.nombre;
  emailInput.value = client.email;
  telefonoInput.value = client.telefono || "";
  direccionInput.value = client.direccion || "";

  formTitle.textContent = "Editar Cliente";
  btnSubmit.textContent = "Actualizar";
  btnCancel.classList.remove("hidden");
  window.scrollTo({ top: 0, behavior: "smooth" });
}

function resetForm() {
  editingId = null;
  form.reset();
  clientIdInput.value = "";
  formTitle.textContent = "Nuevo Cliente";
  btnSubmit.textContent = "Crear";
  btnCancel.classList.add("hidden");
}

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  const payload = {
    nombre: nombreInput.value.trim(),
    email: emailInput.value.trim(),
    telefono: telefonoInput.value.trim(),
    direccion: direccionInput.value.trim(),
  };

  try {
    let res;
    if (editingId) {
      res = await fetch(`${API_BASE}/${editingId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
    } else {
      res = await fetch(API_BASE, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
    }

    if (!res.ok) {
      const err = await res.json();
      throw new Error(err.error || "Error al guardar");
    }

    resetForm();
    await fetchClients();
  } catch (err) {
    showError(err.message);
  }
});

btnCancel.addEventListener("click", resetForm);
btnRefresh.addEventListener("click", fetchClients);

async function deleteClient(id) {
  if (!confirm("¿Estás seguro de eliminar este cliente?")) return;
  try {
    const res = await fetch(`${API_BASE}/${id}`, { method: "DELETE" });
    if (!res.ok && res.status !== 204) {
      const err = await res.json();
      throw new Error(err.error || "Error al eliminar");
    }
    await fetchClients();
  } catch (err) {
    showError(err.message);
  }
}

function showLoading(show) {
  loadingEl.classList.toggle("hidden", !show);
}

function showError(msg) {
  errorEl.textContent = msg;
  errorEl.classList.remove("hidden");
}

function hideError() {
  errorEl.classList.add("hidden");
}

function escapeHtml(str) {
  if (!str) return "";
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

// Init
fetchClients();
