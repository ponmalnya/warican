"use client";

import { useEffect, useState } from "react";

// 疎通確認用。ブラウザ（クライアント）から Go API を直接叩く経路を試す。
// ここで使う URL はホスト側のポートなので NEXT_PUBLIC_API_URL。
// クロスオリジンになるため、api 側の CORS 設定が効いていないと失敗する。
export default function HealthCheck() {
  const [result, setResult] = useState("checking...");
  const apiUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

  useEffect(() => {
    fetch(`${apiUrl}/healthz`)
      .then(async (res) => `${res.status} ${await res.text()}`)
      .then(setResult)
      .catch((err: unknown) => setResult(`NG: ${String(err)}`));
  }, [apiUrl]);

  return (
    <dl>
      <dt>ブラウザ → api（NEXT_PUBLIC_API_URL）</dt>
      <dd>
        <code>{apiUrl}/healthz</code>
      </dd>
      <dd>
        <code>{result}</code>
      </dd>
    </dl>
  );
}
