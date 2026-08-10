function renderLayout(title, content) {
    document.body.innerHTML = `
    <nav class="navbar navbar-dark bg-dark">
        <div class="container-fluid">
            <span class="navbar-brand">
                Go E-Commerce Admin
            </span>

            <button
                id="logout-btn"
                class="btn btn-outline-light">
                Logout
            </button>
        </div>
    </nav>

    <div class="container-fluid">

        <div class="row">

            <div class="col-md-2 sidebar">

                <a href="/dashboard">Dashboard</a>
                <a href="/categories">Categories</a>
                <a href="/products">Products</a>
                <a href="/cart">Cart</a>
                <a href="/orders">Orders</a>

            </div>

            <div class="col-md-10 p-4">

                <h2>${title}</h2>

                <hr>

                ${content}

            </div>

        </div>

    </div>
    `;

    const logoutBtn = document.getElementById("logout-btn");

    if (logoutBtn) {
        logoutBtn.addEventListener("click", logout);
    }
}