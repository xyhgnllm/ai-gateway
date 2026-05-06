import { AuthProvider } from "@/components/auth-provider";
import "./globals.css";

export const metadata = {
  title: "AI Gateway Console",
  description: "AI Gateway user console and admin panel",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
