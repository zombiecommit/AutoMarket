// BLOQUE: utilidades compartidas por las paginas de AutoMarket.
// El token y el usuario los guarda login.html en sessionStorage.
const AutoMarket = (() => {
  const RUTA_FOTO_VACIA = "/static/images/sin-foto.svg";

  function token() { return sessionStorage.getItem("token"); }

  function usuario() {
    try { return JSON.parse(sessionStorage.getItem("usuario")); }
    catch { return null; }
  }

  function rol() {
    const u = usuario();
    return token() && u && u.rol ? u.rol : "visitante";
  }

  function cerrarSesion() {
    sessionStorage.removeItem("token");
    sessionStorage.removeItem("usuario");
    window.location.href = "/catalogo";
  }

  // Menú por rol. Refleja services/menu_service.go, pero solo con las
  // páginas que existen hoy; el servidor no conoce el rol al servir la
  // página (el token vive en el navegador), así que el menú se repinta aquí.
  const MENUS = {
    visitante: [
      ["Catálogo", "/catalogo"], ["Iniciar sesión", "/login"], ["Registrarse", "/registro"],
    ],
    vendedor: [
      ["Catálogo", "/catalogo"], ["Mis vehículos", "/mis-vehiculos"], ["Publicar vehículo", "/publicar"],
    ],
    administrador: [["Catálogo", "/catalogo"]],
  };

  function pintarMenu() {
    const nav = document.querySelector("nav.menu");
    if (!nav) return;

    const lista = nav.querySelector("ul");
    lista.replaceChildren();

    for (const [texto, ruta] of MENUS[rol()] || MENUS.visitante) {
      const li = document.createElement("li");
      const a = document.createElement("a");
      a.href = ruta;
      a.textContent = texto;
      li.appendChild(a);
      lista.appendChild(li);
    }

    const etiquetaRol = nav.querySelector(".rol");
    if (etiquetaRol) etiquetaRol.textContent = "Rol: " + rol();

    if (token()) {
      const boton = document.createElement("button");
      boton.type = "button";
      boton.className = "salir";
      boton.textContent = "Cerrar sesión";
      boton.addEventListener("click", cerrarSesion);
      nav.appendChild(boton);
    }
  }

  // Llama a la API con el token. Si la sesión no es válida, vuelve al login.
  async function api(ruta, opciones = {}) {
    const cabeceras = { Accept: "application/json", ...(opciones.headers || {}) };
    // Con FormData el navegador pone solo el Content-Type (con su boundary).
    if (opciones.body && !(opciones.body instanceof FormData)) {
      cabeceras["Content-Type"] = "application/json";
    }
    if (token()) cabeceras.Authorization = "Bearer " + token();

    const respuesta = await fetch(ruta, { ...opciones, headers: cabeceras });
    const datos = await respuesta.json().catch(() => ({}));

    if (respuesta.status === 401 && opciones.requiereSesion) {
      cerrarSesionSilenciosa();
      window.location.href = "/login";
    }

    return { ok: respuesta.ok, status: respuesta.status, datos };
  }

  function cerrarSesionSilenciosa() {
    sessionStorage.removeItem("token");
    sessionStorage.removeItem("usuario");
  }

  function precio(valor) {
    return new Intl.NumberFormat("es-CO", {
      style: "currency", currency: "COP", maximumFractionDigits: 0,
    }).format(valor);
  }

  // Foto subida al publicar (vehiculo.foto); si no tiene o falla, imagen genérica.
  function foto(vehiculo, texto) {
    const img = document.createElement("img");
    img.alt = texto || `${vehiculo.marca} ${vehiculo.modelo}`;
    img.loading = "lazy";
    img.src = vehiculo.foto || RUTA_FOTO_VACIA;
    img.addEventListener("error", () => {
      img.onerror = null;
      img.src = RUTA_FOTO_VACIA;
    }, { once: true });
    return img;
  }

  // Crea un elemento con texto seguro (nunca interpreta HTML de la API).
  function el(etiqueta, clase, texto) {
    const nodo = document.createElement(etiqueta);
    if (clase) nodo.className = clase;
    if (texto !== undefined) nodo.textContent = texto;
    return nodo;
  }

  const ETIQUETAS_ESTADO = {
    pendiente_aprobacion: "Pendiente de aprobación",
    publicado: "Publicado",
    vendido: "Vendido",
  };

  document.addEventListener("DOMContentLoaded", pintarMenu);

  return { token, usuario, rol, api, precio, foto, el, ETIQUETAS_ESTADO, cerrarSesion };
})();
