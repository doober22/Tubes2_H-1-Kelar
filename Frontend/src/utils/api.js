export async function searchRecipe({ target, method, mode, limit }) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/search`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ target, method, mode, limit }),
  });

  const text = await res.text();
  console.log("Raw response text:", text);

  if (!res.ok) {
    console.error("Server response error:", text);
    throw new Error("Failed to fetch recipe");
  }

  return JSON.parse(text);
}
