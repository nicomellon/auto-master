async function sendData(data) {
	const formData = new FormData();
	const selection = await window.showOpenFilePicker();
	if (selection.length > 0) {
		const file = await selection[0].getFile();
		formData.append("file", file);
	}

	try {
		await fetch("http://localhost:8080/upload", {
			method: "POST",
			body: formData,
		});
	} catch (e) {
		console.error(e);
	}
}

const send = document.querySelector("#upload");
send.addEventListener("click", sendData);

