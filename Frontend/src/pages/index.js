import { useState } from "react";
import RecipeForm from "@/components/RecipeForm";
import StatsPanel from "@/components/StatsPanel";
import RecipeTreeVisual from "@/components/RecipeTreeVisualClient";
import { motion, AnimatePresence } from "framer-motion";

export default function Home() {
  const [result, setResult] = useState(null);

  const handleResult = (data) => {
    setResult(null);      // reset dulu supaya tree lama gak tersisa
    setResult(data);
  };

  return (
    <main className="p-6 space-y-6">
    <h1 className="text-3xl font-bold">Little Alchemy 2</h1>
    <p className="text-lg text-gray-600 mb-4">Element Recipe Search</p>


      <RecipeForm onResult={handleResult} />
        <AnimatePresence>
          {result?.trees?.length > 0 && (
            <>
              <StatsPanel 
              className="sticky top-4 z-10 bg-white p-4 shadow-md rounded-lg"
              time={result.timeMs} 
              nodes={result.nodesVisited} />

              <div className="grid md:grid-cols-2 gap-6">
                {result.trees.map((tree, i) => (
                  <motion.div
                    key={i}
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.4 }}
                    className="p-4 border rounded-xl shadow-sm bg-white"
                  >
                    {result.trees.length > 1 && (
                      <h3 className="font-semibold mb-2 text-blue-600">
                        Recipe #{i + 1}
                      </h3>
                    )}
                    <RecipeTreeVisual data={tree} />
                  </motion.div>
                ))}
              </div>
            </>
          )}
        </AnimatePresence>
    </main>
  );
}
