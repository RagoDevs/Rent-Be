<script lang="ts">
  import type { PageData } from "./$types";
  export let data: PageData;
</script>

<h1>Payment</h1>

<div>
  <p>
    Location: <strong>{data.payment.location}</strong>, Block: {data.payment
      .block}, Partition: {data.payment.partition},
  </p>
  <p>
    Payment ID: {data.payment.id}<br />
    Tenant Name: {data.payment.tenant_name}<br />
    Tenant ID: {data.payment.tenant_id}<br />
    Amount: {data.payment.amount}<br />
    Start Date: {data.payment.start_date}<br />
    End Date: {data.payment.end_date}<br />
    Admin Phone: {data.payment.admin_phone}<br />
    Created At: {data.payment.created_at}<br />
    Updated At: {data.payment.updated_at}<br />
    Version: {data.payment.version}<br />
  </p>
</div>
