document.addEventListener('DOMContentLoaded', function () {
    const dropdownToggle = document.querySelector('.dropdown-toggle');
    const dropdown = document.querySelector('.dropdown');

    if (dropdownToggle && dropdown) {
        dropdownToggle.addEventListener('click', function (event) {
            event.preventDefault();
            dropdown.classList.toggle('show');
        });

        // Close the dropdown if clicked outside
        document.addEventListener('click', function (event) {
            if (!dropdown.contains(event.target)) {
                dropdown.classList.remove('show');
            }
        });
    }
});
