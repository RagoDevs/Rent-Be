<script lang="ts">
  import type { PageData } from "./$types";
  export let data: PageData;
</script>

<h1>Welcome to Tenants</h1>

<div>
  <a href={`/auth/tenants/${data.tenant.id}`}
    ><p>
      First Name: <strong>{data.tenant.first_name}</strong>, Last Name: {data
        .tenant.last_name}, Phone: {data.tenant.phone}
    </p></a
  >
</div>
