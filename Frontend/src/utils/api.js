// src/utils/api.js
export async function searchRecipe({ target, method, mode, limit = 3 }) {
  try {
    const res = await fetch("http://localhost:8080/api/search", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ target, method, mode, limit }),
    });

    if (!res.ok) {
      throw new Error("Search failed");
    }

    const data = await res.json();
    return data;
  } catch (error) {
    console.error("Error in searchRecipe:", error);
    return { found: false, results: [] };
  }
}
