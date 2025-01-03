import { clsx } from "clsx";
import { twMerge } from "tailwind-merge"
import animationData from "../assets/bouncing-dots.json";

export function cn(...inputs) {
  return twMerge(clsx(inputs));
}

export const animationDefaultOptions = {
  loop: true,
  autplay: true,
  animationData, 
};