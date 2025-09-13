import { create } from 'zustand';

export const useStore = create((set, get) => ({
    connectionAlive: null,
    setConnectionAlive: (connectionState) => set({ sharedData: connectionState }),
    allowUploads: null,
    setAllowUploads: (allowState) => set({ allowUploads: allowState }),
    isProcessing: false,
    setIsProcessing: (processingState) => set({ isProcessing: processingState }),
    chatHistory: [],
    setChatHistory: (message) => {
      const currentChatHistory = get().chatHistory;
      set({ chatHistory: [...currentChatHistory, message] });
    },
    resetChatHistory: () => {
      set({ chatHistory: [] });
    },
    showSuggestions: false,
    setShowSuggestions: (state) => set({ showSuggestions: state }),
  }));

export const useWebSocketStore = create((set) => ({
    ws: null,
    connected: null,
    connID: null,
    docCount: 0,
    setWebSocket: (ws) => set({ ws }),
    setConnected: (status) => set({ connected: status }),
    setConnID: (receivedID) => {
      console.log("Setting connID in store:", receivedID);
      set({ connID: receivedID });
    },
    setDocCount: (newCount) => set({docCount: newCount}),
}));