"use client";

import { clearToken, getToken, setToken } from "@/lib/api";
import type { User } from "@/types/api";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

type AuthContextValue = {
  user: User | null;
  token: string | null;
  ready: boolean;
  signIn: (token: string, user: User) => void;
  signOut: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = useState<string | null>(null);
  const [user, setUser] = useState<User | null>(null);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    setTokenState(getToken());
    const rawUser = localStorage.getItem("ai_gateway_user");
    if (rawUser) {
      try {
        setUser(JSON.parse(rawUser));
      } catch {
        localStorage.removeItem("ai_gateway_user");
      }
    }
    setReady(true);
  }, []);

  const signIn = useCallback((nextToken: string, nextUser: User) => {
    setToken(nextToken);
    localStorage.setItem("ai_gateway_user", JSON.stringify(nextUser));
    setTokenState(nextToken);
    setUser(nextUser);
  }, []);

  const signOut = useCallback(() => {
    clearToken();
    localStorage.removeItem("ai_gateway_user");
    setTokenState(null);
    setUser(null);
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      ready,
      signIn,
      signOut,
    }),
    [ready, signIn, signOut, token, user],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) throw new Error("useAuth must be used inside AuthProvider");
  return value;
}
