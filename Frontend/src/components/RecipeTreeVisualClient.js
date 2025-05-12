import dynamic from "next/dynamic";

// Ini yang dipakai di page, bukan RecipeTreeVisual langsung
const RecipeTreeVisualClient = dynamic(
  () => import("./RecipeTreeVisual"),
  { ssr: false }
);

export default RecipeTreeVisualClient;