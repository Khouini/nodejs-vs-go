const express = require('express');
const axios = require('axios');
const app = express();
const port = 3000;

// Middleware to parse JSON requests
app.use(express.json());

// Define the endpoint
app.get('/', async (req, res) => {
    const nbRequests = parseInt(req.query.nbRequests, 10);

    if (isNaN(nbRequests) || nbRequests <= 0) {
        return res.status(400).send('Invalid number of requests. Please provide a positive integer.');
    }

    const startTime = Date.now();

    // Generate an array of URLs based on nbRequests
    const urls = Array.from({length: nbRequests}, (_, i) => `https://jsonplaceholder.typicode.com/photos/${i + 1}`);

    try {
        // Create an array to hold the promises and processing times
        const promises = [];
        const processingTimes = [];

        urls.forEach(url => {
            const requestStartTime = Date.now();
            const promise = axios.get(url).then(response => {
                const requestEndTime = Date.now();
                processingTimes.push({
                    url,
                    processingTime: requestEndTime - requestStartTime
                });
                return response;
            });
            promises.push(promise);
        });

        // Wait for all promises to resolve
        const responses = await Promise.all(promises);

        // Process the responses
        const results = responses.map(response => response.data);

        // Calculate the maximum processing time
        const maxProcessingTime = Math.max(...processingTimes.map(pt => pt.processingTime));

        res.status(200).json({
            totalProcessingTime: Date.now() - startTime,
            maxProcessingTime,
            // processingTimes,
            // results,
        });
    } catch (error) {
        console.error('Error fetching data:', error);
        res.status(500).send('Error fetching data');
    }
});

app.get("/heavy", (req, res) => {
    const startTime = Date.now();
    const iterations = 1000000000;
    let result = 0;
    for (let i = 0; i < iterations; i++) {
        result += i;
    }
    const endTime = Date.now();
    res.status(200).json({
        totalProcessingTime: endTime - startTime,
        result
    });
});

app.get("/hello", (req, res) => {
    res.status(200).send("Hello World!");
})
    // heavy computation

// Start the server
app.listen(port, () => {
    console.log(`Server running at http://localhost:${port}`);
});
