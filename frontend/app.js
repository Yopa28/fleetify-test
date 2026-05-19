const API_URL = "http://localhost:8080/api";

const vehicleSelect = document.getElementById("vehicle");
const reportTableBody = document.querySelector("#reportTable tbody");
const reportForm = document.getElementById("reportForm");
const exportCsvBtn = document.getElementById("exportCsvBtn");

function clearElement(element) {
    while (element.firstChild) {
        element.removeChild(element.firstChild);
    }
}

async function fetchVehicles() {
    try {
        const response = await fetch(`${API_URL}/vehicles`);
        const result = await response.json();

        clearElement(vehicleSelect);

        result.data.forEach((vehicle) => {
            const option = document.createElement("option");

            option.value = vehicle.ID;
            option.textContent =
                `${vehicle.LicensePlate} - ${vehicle.Model}`;

            vehicleSelect.appendChild(option);
        });

    } catch (error) {
        console.error("Failed to fetch vehicles", error);
    }
}

async function fetchReports() {
    try {
        const response = await fetch(`${API_URL}/reports`);
        const result = await response.json();

        clearElement(reportTableBody);

        result.data.forEach((report) => {

            const tr = document.createElement("tr");

            const tdId = document.createElement("td");
            tdId.textContent = report.ID;

            const tdVehicle = document.createElement("td");
            tdVehicle.textContent =
                report.Vehicle.LicensePlate;

            const tdStatus = document.createElement("td");
            tdStatus.textContent = report.Status;

            const tdComplaint = document.createElement("td");
            tdComplaint.textContent = report.Complaint;

            tr.appendChild(tdId);
            tr.appendChild(tdVehicle);
            tr.appendChild(tdStatus);
            tr.appendChild(tdComplaint);

            reportTableBody.appendChild(tr);
        });

    } catch (error) {
        console.error("Failed to fetch reports", error);
    }
}

async function exportCSV() {
    try {
        const response = await fetch(`${API_URL}/reports`);
        const result = await response.json();

        const reports = result.data;

        let csvContent =
            "ID,Vehicle,Status,Complaint,Odometer\n";

        reports.forEach((report) => {
            csvContent +=
                `${report.ID},` +
                `${report.Vehicle.LicensePlate},` +
                `${report.Status},` +
                `"${report.Complaint}",` +
                `${report.Odometer}\n`;
        });

        const blob = new Blob(
            [csvContent],
            { type: "text/csv" }
        );

        const url = window.URL.createObjectURL(blob);

        const link = document.createElement("a");

        link.href = url;
        link.download = "fleetify-reports.csv";

        document.body.appendChild(link);

        link.click();

        document.body.removeChild(link);

        window.URL.revokeObjectURL(url);

    } catch (error) {
        console.error("Failed to export CSV", error);
    }
}

reportForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const payload = {
        vehicle_id: parseInt(vehicleSelect.value),
        odometer: parseInt(
            document.getElementById("odometer").value
        ),
        complaint:
            document.getElementById("complaint").value,
        initial_photo: "frontend-upload.jpg",
        items: [
            {
                item_id: 1,
                quantity: 1,
            },
        ],
    };

    try {
        const response = await fetch(`${API_URL}/reports`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-User-ID": "1",
            },
            body: JSON.stringify(payload),
        });

        const result = await response.json();

        alert(result.message);

        reportForm.reset();

        fetchReports();

    } catch (error) {
        console.error("Failed to create report", error);
    }
});

fetchVehicles();
fetchReports();
exportCsvBtn.addEventListener("click", exportCSV);