import { useState } from "react";
import RecipeForm from "@/components/RecipeForm";
import StatsPanel from "@/components/StatsPanel";
import RecipeTreeVisual from "@/components/RecipeTreeVisualClient";
import { motion, AnimatePresence } from "framer-motion";

export default function Home() {
  const [result, setResult] = useState(null);
  const [version, setVersion] = useState(0);

  const handleResult = (data) => {
    setVersion((v) => v + 1);
    setResult(data);
  };

  return (
    <main className="min-h-screen px-6 py-10 bg-gradient-to-tr from-blue-50 via-purple-50 to-pink-50">
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-extrabold text-purple-700 drop-shadow-sm">
            🧪 Little Alchemy 2
          </h1>
          <p className="text-lg text-gray-700 mt-2">Element Recipe Search</p>
        </div>

        <div className="bg-white p-6 rounded-2xl shadow-md ring-1 ring-gray-200">
          <RecipeForm onResult={handleResult} />
        </div>

        <AnimatePresence mode="wait">
          {result?.trees?.length > 0 && (
            <motion.div
              key={`trees-wrapper-${version}`}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.4 }}
              className="space-y-6"
            >
              <StatsPanel
                className="sticky top-4 z-10 bg-white p-4 shadow-lg rounded-xl border border-purple-100"
                time={result.timeMs}
                nodes={result.nodesVisited}
              />

              
                {result.trees.map((tree, i) => (
                  <motion.div
                    key={`tree-${version}-${i}`}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    transition={{ duration: 0.4 }}
                    className="p-5 border border-gray-200 rounded-xl shadow-sm bg-white hover:shadow-md transition-shadow"
                  >
                    {result.trees.length > 1 && (
                      <h3 className="font-semibold mb-3 text-indigo-600 text-lg">
                        Recipe #{i + 1}
                      </h3>
                    )}
                    <RecipeTreeVisual data={tree} />
                  </motion.div>
                ))}
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </main>
  );
}
