// src/components/RecipeTreeVisual.js
import { Tree, TreeNode } from 'react-organizational-chart';

function buildTree(steps) {
  const nodeMap = new Map();
  let root = null;

  for (const step of steps) {
    const parts = step.replace(/ /g, '').split(/\+|=/);
    if (parts.length !== 3) continue;
    const [left, right, result] = parts;

    let leftNode = nodeMap.get(left);
    if (!leftNode) {
      leftNode = { name: left, children: [] };
      nodeMap.set(left, leftNode);
    }

    let rightNode = nodeMap.get(right);
    if (!rightNode) {
      rightNode = { name: right, children: [] };
      nodeMap.set(right, rightNode);
    }

    let resultNode = nodeMap.get(result);
    if (!resultNode) {
      resultNode = { name: result, children: [] };
      nodeMap.set(result, resultNode);
    }

    // Hindari overwrite jika children sudah terisi
    if (resultNode.children.length === 0) {
      resultNode.children.push(leftNode, rightNode);
    }

    root = resultNode;
  }

  return root;
}

export default function RecipeTreeVisual({ steps }) {
  if (!steps || steps.length === 0) return null;

  const tree = buildTree(steps);

  const renderTree = (node) => (
    <TreeNode key={node.name} label={<div className="border rounded px-2 py-1 bg-white shadow text-sm">{node.name}</div>}>
      {node.children.map((child) => renderTree(child))}
    </TreeNode>
  );

  return (
    <div className="overflow-auto p-4 border rounded">
      <h2 className="font-semibold mb-2">Visual Recipe Tree</h2>
      <Tree
        lineWidth={'2px'}
        lineColor={'#ccc'}
        lineBorderRadius={'8px'}
        label={<div className="border rounded px-2 py-1 bg-blue-100 shadow text-sm font-bold">{tree.name}</div>}
      >
        {tree.children.map((child) => renderTree(child))}
      </Tree>
    </div>
  );
}