import { Tree, TreeNode } from "react-organizational-chart";

export default function RecipeTreeVisual({ data }) {
  if (!data) return null;

  const renderTree = (node) => (
    <TreeNode
      key={node.element}
      label={
        <div className="border rounded px-2 py-1 bg-white shadow text-sm">
          {node.element}
        </div>
      }
    >
      {node.ingredients?.map((child) => renderTree(child))}
    </TreeNode>
  );

  return (
    <div className="overflow-auto p-4 border rounded">
      <h2 className="font-semibold mb-2">Visual Recipe Tree</h2>
      <Tree
        lineWidth={"2px"}
        lineColor={"#ccc"}
        lineBorderRadius={"8px"}
        label={
          <div className="border rounded px-2 py-1 bg-blue-100 shadow text-sm font-bold">
            {data.element}
          </div>
        }
      >
        {data.ingredients?.map((child) => renderTree(child))}
      </Tree>
    </div>
  );
}