import { useState } from "react";
import RecipeForm from "@/components/RecipeForm";
import StatsPanel from "@/components/StatsPanel";
import RecipeTreeVisual from "@/components/RecipeTreeVisualClient";

export default function Home() {
  const [result, setResult] = useState(null);

  const handleResult = (data) => {
    setResult(null);      // reset dulu supaya tree lama gak tersisa
    setResult(data);
  };

  return (
    <main className="p-6 space-y-6">
      <h1 className="text-2xl font-bold">Little Alchemy 2</h1>
      <h1 className="text-2xl font-bold">Element Recipe Search</h1>

      <RecipeForm onResult={handleResult} />

      {result?.trees?.length > 0 && (
        <>
          <StatsPanel
            time={result.timeMs}
            nodes={result.nodesVisited}
          />

          {result.trees.map((tree, i) => (
            <div key={i}>
              {result.trees.length > 1 && (
                <h3 className="font-semibold mt-4">Recipe #{i + 1}</h3>
              )}
              <RecipeTreeVisual data={tree} />
            </div>
          ))}
        </>
      )}
    </main>
  );
}
