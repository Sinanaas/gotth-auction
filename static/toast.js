class Toast {
    constructor(level, message) {
        this.level = level;
        this.message = message;
    }

    #makeToastContainerButton() {
        const button = document.createElement("button");
        button.classList.add("toast", `toast-${this.level}`, "px-4", "py-2", "rounded", "m-2");
        button.setAttribute("role", "alert");
        button.setAttribute("aria-label", "Close");
        button.addEventListener("click", () => button.remove());
        return button;
    }

    #makeToastContentElement() {
        const messageContainer = document.createElement("span");
        messageContainer.textContent = this.message;
        return messageContainer;
    }

    show(containerQuerySelector = "#toast-container") {
        const toast = this.#makeToastContainerButton();
        const toastContent = this.#makeToastContentElement();
        toast.appendChild(toastContent);

        const toastContainer = document.querySelector(containerQuerySelector);
        toastContainer.appendChild(toast);
    }
}

document.body.addEventListener("makeToast", onMakeToast);

function onMakeToast(e) {
    const toast = new Toast(e.detail.level, e.detail.message);
    toast.show();
}
