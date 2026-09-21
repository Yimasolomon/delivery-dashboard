console.log("Delivery Dashboard loaded.");

const chartElement = document.getElementById("statusChart");

if (chartElement && typeof Chart !== "undefined") {
    const labels = statusCounts.map(item => item.status);
    const data = statusCounts.map(item => item.count);

    new Chart(chartElement, {
        type: "doughnut",
        data: {
            labels: labels,
            datasets: [
                {
                    data: data
                }
            ]
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    position: "bottom"
                }
            }
        }
    });
}


const trendChartElement = document.getElementById("deliveryTrendChart");

if (trendChartElement && typeof Chart !== "undefined") {
    const labels = deliveryTrends.map(item => item.date);
    const data = deliveryTrends.map(item => item.count);

    new Chart(trendChartElement, {
        type: "line",
        data: {
            labels: labels,
            datasets: [
                {
                    label: "Deliveries",
                    data: data,
                    fill: false,
                    tension: 0.3
                }
            ]
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    position: "bottom"
                }
            }
        }
    });
}