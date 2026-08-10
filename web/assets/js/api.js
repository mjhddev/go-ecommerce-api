const API_BASE_URL = "/api/v1";

async function apiRequest(endpoint, options = {}) {
    const token = getToken();

    const headers = {
        "Content-Type": "application/json",
        ...options.headers,
    };

    if (token) {
        headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        ...options,
        headers,
    });

    const result = await response.json();

    if (!response.ok) {
        throw new Error(result.message || "Request failed");
    }

    return result;
}