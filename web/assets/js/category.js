protectPage();

document.addEventListener("DOMContentLoaded", () => {
    loadCategories();
});

async function loadCategories() {
    try {

        const result = await apiRequest("/categories");

        renderCategories(result.data);

    } catch (err) {

        console.error(err);

    }
}

function renderCategories(categories) {

    const tbody = document.getElementById("category-table");

    tbody.innerHTML = "";

    categories.forEach(category => {

        tbody.innerHTML += `
            <tr>

                <td>${category.id}</td>

                <td>${category.name}</td>

                <td>

                    <button
                        class="btn btn-warning btn-sm">

                        Edit

                    </button>

                    <button
                        class="btn btn-danger btn-sm">

                        Delete

                    </button>

                </td>

            </tr>
        `;

    });

}