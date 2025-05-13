// pages/api/search.js

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const { target, method, mode, limit } = req.body;

  try {
    const backendRes = await fetch("http://localhost:8080/search", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ target, method, mode, limit }),
    });

    if (!backendRes.ok) {
      const errorText = await backendRes.text();
      return res.status(backendRes.status).json({ error: errorText });
    }

    const data = await backendRes.json();
    return res.status(200).json(data);
  } catch (error) {
    console.error("Error contacting Go backend:", error);
    return res.status(500).json({ error: "Internal server error" });
  }
}
