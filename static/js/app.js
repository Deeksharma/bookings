// Prompt is our JS module for alerts, notifications and custom page dialogs
function Prompt() {
    let toast = function (c) {
        const {
            msg = "", // these are default values
            icon = "success",
            position = "top-end",
        } = c; // storing the variable without name in c, c will be overwritten by the things that are here
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }
    let success = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c;
        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
        })
    }

    let error = function (c) { // these are modules
        const {
            msg = "",
            title = "",
            footer = "",
        } = c;
        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
        })
    }

    let custom = async function (c) {
        const {
            icon = "",
            msg = "",
            title = "",
            callback = function () {
                console.log("default callback function")
            },
            willOpen = undefined,
            didOpen = undefined,
            showConfirmButton = true,
        } = c;

        const {value: result} = await Swal.fire({ // can be used only in async method
            icon: icon,
            title: title,
            html: msg, // html code
            focusConfirm: false,
            backdrop: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            }
        })

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) // if didnt hit the cancel button
            {
                if (result.value !== "") { // I want to execute a JS back on the client page which happens when client put the date
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

(function () {
    'use strict'

    // Fetch all the forms we want to apply custom Bootstrap validation styles to
    let forms = document.querySelectorAll('.needs-validation')

    // Loop over them and prevent submission
    Array.prototype.slice.call(forms)
        .forEach(function (form) {
            form.addEventListener('submit', function (event) {
                if (!form.checkValidity()) {
                    event.preventDefault()
                    event.stopPropagation()
                }

                form.classList.add('was-validated')
            }, false)
        })
})()
