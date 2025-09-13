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

export const isIOSSafari = () => {
  const userAgent = window.navigator.userAgent;
  const isIOS = /iPad|iPhone|iPod/.test(userAgent) && !window.MSStream;
  const isSafari = /Safari/.test(userAgent) && !/Chrome/.test(userAgent);
  return isIOS && isSafari;
};

export const isIDPreset = (id) => {
  if (id.startsWith("@")) {
    return true;
  }
  return false;
}