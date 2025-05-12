// pages/api/search.js
export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const { target } = req.body;

  // Delay simulasi
  await new Promise((r) => setTimeout(r, 500));

  const tree = {
    name: target || 'Unknown',
    children: [
      {
        name: 'ElementA',
        children: [
          { name: 'ElementB', children: [] },
          { name: 'ElementC', children: [] },
        ],
      },
      { name: 'ElementD', children: [] },
    ],
  };

  res.status(200).json({
    tree,
    time: 42,
    nodesVisited: 13,
  });
}
