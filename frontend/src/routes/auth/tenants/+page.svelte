<script lang="ts">
  import type { PageData } from "./$types";
  export let data: PageData;
</script>

<h1>Welcome to Tenants</h1>

{#each data.tenants as tenant}
  <div>
    <a href={`/auth/tenants/${tenant.id}`}
      ><p>
       First Name: <strong>{tenant.first_name}</strong>, Last Name: {tenant.last_name}, Phone: {tenant.phone}
      </p></a
    >
  </div>
{/each}



