import dynamic from "next/dynamic";

const RecipeTreeVisualClient = dynamic(
  () => import("./RecipeTreeVisual"),
  { ssr: false }
);

export default RecipeTreeVisualClient;