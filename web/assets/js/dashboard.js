protectPage();

document.addEventListener("DOMContentLoaded", async () => {

    renderLayout(
        "Dashboard",
        `
        <div class="row">

            <div class="col-md-3">
                <div class="card shadow-sm">
                    <div class="card-body text-center">
                        <h6>Categories</h6>
                        <h2 id="total-category">0</h2>
                    </div>
                </div>
            </div>

            <div class="col-md-3">
                <div class="card shadow-sm">
                    <div class="card-body text-center">
                        <h6>Products</h6>
                        <h2 id="total-product">0</h2>
                    </div>
                </div>
            </div>

            <div class="col-md-3">
                <div class="card shadow-sm">
                    <div class="card-body text-center">
                        <h6>Cart</h6>
                        <h2 id="total-cart">0</h2>
                    </div>
                </div>
            </div>

            <div class="col-md-3">
                <div class="card shadow-sm">
                    <div class="card-body text-center">
                        <h6>Orders</h6>
                        <h2 id="total-order">0</h2>
                    </div>
                </div>
            </div>

        </div>
        `
    );

    await loadDashboard();
});

async function loadDashboard() {

    try {

        const categories = await apiRequest("/categories");
        const products = await apiRequest("/products");
        const cart = await apiRequest("/cart");
        const orders = await apiRequest("/orders");

        document.getElementById("total-category").innerText =
            categories.data.length;

        document.getElementById("total-product").innerText =
            products.data.length;

        document.getElementById("total-cart").innerText =
            cart.data.length;

        document.getElementById("total-order").innerText =
            orders.data.length;

    } catch (err) {

        console.error(err);

    }

}