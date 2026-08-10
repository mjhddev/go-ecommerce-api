// Login
async function login(event) {
    event.preventDefault();

    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    const errorElement = document.getElementById("error-message");

    // Reset error
    if (errorElement) {
        errorElement.classList.add("d-none");
        errorElement.innerText = "";
    }

    try {
        const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email,
                password,
            }),
        });

        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.message || "Login failed");
        }

        // Simpan JWT
        localStorage.setItem("token", result.data.access_token);

        // Pindah ke dashboard
        window.location.href = "/dashboard";
    } catch (err) {
        if (errorElement) {
            errorElement.classList.remove("d-none");
            errorElement.innerText = err.message;
        }
    }
}

// Cek apakah user sudah login
function protectPage() {
    const token = localStorage.getItem("token");

    if (!token) {
        window.location.href = "/";
    }
}

// Ambil token
function getToken() {
    return localStorage.getItem("token");
}

// Logout
function logout() {
    localStorage.removeItem("token");
    window.location.href = "/";
}

// Pasang event login jika ada form login
document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("login-form");

    if (loginForm) {
        loginForm.addEventListener("submit", login);
    }

    const logoutBtn = document.getElementById("logout-btn");

    if (logoutBtn) {
        logoutBtn.addEventListener("click", logout);
    }
});