function copyText(text) {
	navigator.clipboard.writeText(text);
}

function copyInviteLink(guid) {
	copyText(window.location.origin + '/register?guid=' + guid);
}

function filterAdminTable(input) {
	const section = input.closest('.admin-section');
	if (!section) return;
	const query = input.value.toLowerCase().trim();
	section.querySelectorAll('tbody tr').forEach(function(row) {
		row.style.display = row.textContent.toLowerCase().includes(query) ? '' : 'none';
	});
}

window.filterAdminTable = filterAdminTable;

function closeAdminModal() {
	const host = document.getElementById('modal-host');
	if (host) host.innerHTML = '';
}

window.closeAdminModal = closeAdminModal;
document.addEventListener('closeModal', closeAdminModal);

document.addEventListener('htmx:beforeSwap', function(evt) {
	const xhr = evt.detail.xhr;
	if (!xhr || xhr.status < 200 || xhr.status >= 300) return;
	if (xhr.responseText && xhr.responseText.trim() !== '') return;
	const target = evt.detail.target;
	if (target && target.tagName === 'TR') {
		evt.detail.shouldSwap = false;
		target.classList.add('htmx-swapping');
		setTimeout(function() { target.remove(); }, 300);
	}
});
