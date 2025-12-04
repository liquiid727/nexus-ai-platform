const express = require("express");
const cors = require("cors");
const helmet = require("helmet");
const swaggerUi = require("swagger-ui-express");
const { createProxyMiddleware } = require("http-proxy-middleware");
const fs = require("fs");
const path = require("path");

const app = express();
app.use(helmet());
app.use(cors());
app.use(express.json());

const openapiPath = path.join(__dirname, "openapi.json");
const openapi = JSON.parse(fs.readFileSync(openapiPath, "utf8"));
app.use("/docs", swaggerUi.serve, swaggerUi.setup(openapi));

app.use(
  "/auth",
  createProxyMiddleware({ target: "http://localhost:3001", changeOrigin: true })
);

app.get('/admin', (req, res) => {
  res.sendFile(path.join(__dirname, 'admin.html'));
});

app.get("/health", (req, res) => {
  res.status(200).send({ status: "ok" });
});

app.get("/version", (req, res) => {
  res.send({ version: "0.1.0" });
});

const port = process.env.PORT || 3000;
app.listen(port, () => {
  console.log(`Gateway listening at http://localhost:${port}`);
});
