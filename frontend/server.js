const express = require("express");
const path = require("path");
const cors = require("cors");

const app = express();

app.use(cors())
app.use("/static", express.static(path.resolve(__dirname, "static")));

app.get("/*", (req, res) => {
    res.sendFile(path.resolve(__dirname, "index.html"));
});

app.listen(process.env.FRONTEND_PORT || 8080, () => console.log("Server running..."));