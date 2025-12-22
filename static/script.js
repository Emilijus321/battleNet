// Toast functionality for GoLand project
console.log("Toast.js loaded in GoLand project");

const ToastManager = {
    init() {
        this.autoRemoveExisting();
        this.setupHtmxListeners();
    },

    autoRemoveExisting() {
        setTimeout(() => {
            document.querySelectorAll('.toast').forEach(toast => {
                toast.remove();
            });
        }, 5000);
    },

    setupHtmxListeners() {
        document.addEventListener('htmx:afterSwap', () => {
            setTimeout(() => {
                document.querySelectorAll('.toast').forEach(toast => {
                    if (!toast.dataset.autoRemoved) {
                        toast.dataset.autoRemoved = 'true';
                        setTimeout(() => toast.remove(), 5000);
                    }
                });
            }, 100);
        });
    },

    show(message, type = 'success') {
        const container = document.getElementById('toast-container');
        if (!container) return;

        const toast = document.createElement('div');
        toast.className = `toast toast-${type}`;
        toast.innerHTML = `
            ${message}
            <button class="close" onclick="this.parentElement.remove()">×</button>
        `;

        container.appendChild(toast);
        setTimeout(() => toast.remove(), 5000);
    }
};

// Auto-initialize when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => ToastManager.init());
} else {
    ToastManager.init();
}

/*--------------------*/
document.addEventListener('DOMContentLoaded', function() {
    var input = document.querySelector('input[name="q"]');
    if (input) input.focus();

    document.querySelectorAll('.tag').forEach(function(tag) {
        tag.onclick = function(e) {
            e.preventDefault();
            var input = document.querySelector('input[name="q"]');
            if (input) {
                input.value = this.textContent;
                input.form.submit();
            }
        };
    });
});