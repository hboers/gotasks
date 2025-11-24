import { ref } from "vue";
import { fetchMe } from "../src/api";

export const currentUser = ref(null);

export async function loadUser() {
  try {
    currentUser.value = await fetchMe();
  } catch {
    currentUser.value = null;
  }
}