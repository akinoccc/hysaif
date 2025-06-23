<script setup lang="ts">
import { RouterView } from 'vue-router'
import { Toaster } from '@/components/ui/sonner'
import 'vue-sonner/style.css'
</script>

<template>
  <div id="app" class="min-h-screen bg-background">
    <RouterView />
    <Toaster />
  </div>
</template>
