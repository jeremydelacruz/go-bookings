function Prompt() {
    let toast = function (c) {
        const{
            msg = '',
            icon = 'success',
            position = 'top-end',

        } = c;

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
        } = c

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    let error = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            }
        })

        if (result) {
            const resultHasValue = result.dismiss !== Swal.DismissReason.cancel && result.value !== "";
            if (resultHasValue) {
                if (c.callback !== undefined) {
                    await c.callback(result);
                }
            } else {
                await c.callback(false);
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

function renderAvailabilityForm(csrfToken, pageRoomId) {
    document.getElementById("check-availability-button").addEventListener("click", function () {
        let html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
            <div class="form-row">
                <div class="col">
                    <div class="form-row" id="reservation-dates-modal">
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                        </div>
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                        </div>

                    </div>
                </div>
            </div>
        </form>
        `;
        attention.custom({
            msg: html,
            title: 'Choose your dates',
            willOpen: () => {
                const elem = document.getElementById("reservation-dates-modal");
                const rp = new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                    minDate: new Date(),
                })
            },
            didOpen: () => {
                document.getElementById("start").removeAttribute("disabled");
                document.getElementById("end").removeAttribute("disabled");
            },
            callback: async function(result) {
                console.log(`called with result: ${result}`);

                if (!result)
                    return;

                const form = document.getElementById("check-availability-form");
                const formData = new FormData(form);
                formData.append("csrf_token", csrfToken);
                formData.append("room_id", `${pageRoomId}`);

                const res = await fetch('/search-availability-json', {
                    method: "post",
                    body: formData
                });
                const data = await res.json();

                if (!data.ok) {
                    attention.error({
                        msg: 'No availability',
                    });
                    return;
                }

                const htmlStr = `
                    <p>It's available!</p>
                    <p><a href="/book-room?id=${data.room_id}&s=${data.start_date}&e=${data.end_date}" class="btn btn-primary">Book now</a></p>
                `;
                attention.custom({
                    icon: 'success',
                    msg: htmlStr,
                    showConfirmButton: false,
                });
            }
        });
    });
}
