import { SvelteKitAuth } from "@auth/sveltekit";
import OryHydra from "@auth/core/providers/ory-hydra";

let isRefreshing = false;
const client_id = "44e6395d-7b22-4de6-9b73-0cc0242c524d";
const client_secret = "b-q.AHE--GX484xQMPmyOUaqeS";

function generateRandomState(length = 8) {
  return [...Array(length)].map(() => Math.random().toString(36)[2]).join("");
}

const HYDRA_TOKEN_ENDPOINT = "http://127.0.0.1:8080/oauth2/token";
async function refreshAccessToken(refresh_token) {
  const response = await fetch(HYDRA_TOKEN_ENDPOINT, {
      method: "POST",
      headers: {
          "Content-Type": "application/x-www-form-urlencoded",
          "Authorization": `Basic ${btoa(`${client_id}:${client_secret}`)}`,
      },
      body: new URLSearchParams({
          grant_type: "refresh_token",
          refresh_token: refresh_token,
      }),
  });

    const tokensOrError = await response.json();
    if (!response.ok) throw tokensOrError;

    return tokensOrError;
}

export const { handle, signIn, signOut } = SvelteKitAuth({
  providers: [
    OryHydra({
      clientId: "44e6395d-7b22-4de6-9b73-0cc0242c524d",
      clientSecret: "b-q.AHE--GX484xQMPmyOUaqeS",
      issuer: "http://127.0.0.1:8080/",
      authorization: {
        params: {
          scope: "openid offline hydra.openid introspect",
          state: generateRandomState(8),
        },
      },
    }),
  ],
  callbacks: {
    async jwt({ token, account }) {
      if (isRefreshing) return token;
      isRefreshing = true;
      try {
      if (account) {
        console.log('Original TOKEN ======= ', {
          ...token,
          access_token: account.access_token,
          expires_at: account.expires_at,
          refresh_token: account.refresh_token,
        })
        return {
          ...token,
          access_token: account.access_token,
          expires_at: account.expires_at,
          refresh_token: account.refresh_token,
        }
      } else if (Date.now() < (token.expires_at - 0 * 60) * 1000) {
        return token
      } else {
        if (!token.refresh_token) throw new TypeError("Missing refresh_token")
        try {
          const newTokens = await refreshAccessToken(token.refresh_token)
          return {
            ...token,
            access_token: newTokens.access_token,
            expires_at: Math.floor(Date.now() / 1000 + newTokens.expires_in),
            refresh_token: newTokens.refresh_token || token.refresh_token,
            error: null,
          }
        } catch (error) {
          console.error("Error refreshing access_token", error)
          token.error = "RefreshTokenError"
          
          return token
        }
      }
    } finally {
      isRefreshing = false;
    }
    },
    async session({ session, token }) {
      session.access_token = token.access_token;
      session.error = token.error;
      return session;
    },
  }
});