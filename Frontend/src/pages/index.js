import RecipeForm from "@/components/RecipeForm";
import StatsPanel from "@/components/StatsPanel";
import RecipeTreeVisual from "@/components/RecipeTreeVisualClient";
import { useState } from "react";

export default function Home() {
  const [result, setResult] = useState(null);

  return (
    <main className="p-6 space-y-6">
      <h1 className="text-2xl font-bold">Element Recipe Search</h1>
      <RecipeForm onResult={setResult} />

      {result?.results?.length > 0 && (
        <>
          <StatsPanel
            time={result.results[0].timeMs}
            nodes={result.results[0].nodesVisited}
          />
          <RecipeTreeVisual steps={result.results[0].steps} />
        </>
      )}
    </main>
  );
}
