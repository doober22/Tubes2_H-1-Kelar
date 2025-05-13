export default function StatsPanel({ time, nodes }) {
  return (
    <div className="p-4 border rounded bg-gray-50">
      <h2 className="font-semibold mb-2">Search Statistics</h2>
      <p><strong>Execution Time:</strong> {time?.toFixed(3)} ms</p>
      <p><strong>Nodes Visited:</strong> {nodes ?? 0}</p>
    </div>
  );
}