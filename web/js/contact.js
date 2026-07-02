// ── Contact form — Formspree submit ──
document.getElementById('tsunami-form').addEventListener('submit', function(e) {
    e.preventDefault();

    const form = this;
    const formData = new FormData(form);

    fetch('https://formspree.io/f/meenwqzy', {
        method: 'POST',
        body: formData,
        headers: {
            'Accept': 'application/json'
        }
    }).then(response => {
        if (response.ok) {
            // Hide form and show success message
            form.style.display = 'none';
            document.getElementById('successMsg').style.display = 'block';
        } else {
            alert('Something went wrong. Please try again.');
        }
    }).catch(error => {
        alert('Something went wrong. Please try again.');
    });
});
