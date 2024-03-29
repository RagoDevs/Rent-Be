<script lang="ts">
  import type { PageData } from "./$types";
  export let data: PageData;
</script>

<h1>Welcome to Payments</h1>

{#each data.payments as payment}
  <div>
    <a href={`/auth/payments/${payment.id}`}>
      <p>
        Amount: {payment.amount}, Start Date: {payment.start_date}, End Date: {payment.end_date}
      </p>
    </a>
  </div>
{/each}


