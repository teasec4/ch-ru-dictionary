<script lang="ts">
    let searchTerm = $state("");
    let responseData = $state<Response[]>([]);
    
    type Response = {
      chinese: string;
      pinyin: string;
      pininNormalized: string;
      meaning: string;
    }

    async function handleClick() {
      if(searchTerm.trim()){
        responseData = await fetch(`/api/search/${searchTerm}`); 
      }
    }
</script>

<h1>Словарь</h1>
<input type="text" placeholder="Search..." bind:value={searchTerm} />
<button onclick={handleClick}>поиск</button>
