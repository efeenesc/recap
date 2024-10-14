<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  interface MousePosition {
    x: number,
    time: number,
  }

  let carousel: HTMLDivElement;
  let contentDiv: HTMLDivElement;
  let carouselRect: DOMRect;
  let carouselBounds: { maxLeft: number, maxRight: number} = { maxLeft: 0, maxRight: 0 };
  let carouselTween: gsap.core.Tween;
  let isDragging = false;
  let dragStartPos = 0;
  let translatePos = 0;
  let elTranslatePos = {
    current: 0
  };
  let currentMousePos: MousePosition | null = { x: 0, time: 0 };
  let prevMousePos: MousePosition | null = { x: 0, time: 0 };

  function onWindowSizeChange() {
    if (!contentDiv) return;
    carouselBounds = {
      maxLeft: -(contentDiv.clientWidth * 0.85),
      maxRight: 0
    }
  }

  function clearMousePosition() {
    prevMousePos = null;
    currentMousePos = null;
  }

  function getDragPosition(e: MouseEvent | TouchEvent, left: number): number {
    if (e instanceof TouchEvent) {
      return e.changedTouches[0].clientX - left;
    } else {
      return e.clientX - left;
    }
  }

  function setPrevMousePosition(x: number, time: number): void {
    if (currentMousePos) prevMousePos = { ...currentMousePos };
    currentMousePos = { x, time };
  }

  function startDragging(e: MouseEvent | TouchEvent): void {
    if (isDragging) return;
    isDragging = true;
    clearMousePosition();

    // Calculate carousel rect on drag start
    carouselRect = carousel.getBoundingClientRect();

    dragStartPos = getDragPosition(e, carouselRect.left);
    setPrevMousePosition(dragStartPos, e.timeStamp);
  }

  function drag(e: MouseEvent | TouchEvent): void {
    if (!isDragging) return;

    const offsetX = getDragPosition(e, carouselRect.left);
    const dragDistance = offsetX - dragStartPos;

    translatePos += dragDistance;
    dragStartPos = offsetX;

    setPrevMousePosition(offsetX, e.timeStamp);
    updateCarouselPosition(150);
  }

  function stopDragging(): void {
    if (!isDragging) return;
    isDragging = false;

    if (!prevMousePos || !currentMousePos) return;
  
    const { x: curX, time: curTime } = currentMousePos;
    const { x: prevX, time: prevTime } = prevMousePos;
  
    const dt = curTime - prevTime;
    const dx = curX - prevX;
  
    const vel = Math.round((dx / dt) * 100); // Calculate velocity
    translatePos += vel;
    translatePos = Math.max(
      Math.min(translatePos, carouselBounds.maxRight),
      carouselBounds.maxLeft
    );
    updateCarouselPosition(750);
  }

  onMount(() => {
    carouselRect = carousel.getBoundingClientRect();
    window.addEventListener('resize', onWindowSizeChange);

    // Create a ResizeObserver to watch for changes in the content div
    const resizeObserver = new ResizeObserver(() => {
      onWindowSizeChange();
    });

    // Observe the content div
    resizeObserver.observe(contentDiv);

    window.addEventListener('mousemove', drag);
    window.addEventListener('mouseup', stopDragging);
    window.addEventListener('touchmove', drag);
    window.addEventListener('touchend', stopDragging);

    return () => {
      resizeObserver.disconnect();
      window.removeEventListener('resize', onWindowSizeChange);
      window.removeEventListener('mousemove', drag);
      window.removeEventListener('mouseup', stopDragging);
      window.removeEventListener('touchmove', drag);
      window.removeEventListener('touchend', stopDragging);
    }
  });

  function onLoad() {
    onWindowSizeChange();
  }

  function updateCarouselPosition(duration: number) {
    if (carouselTween)
      carouselTween.kill();

    carouselTween = gsap.to(elTranslatePos, {
      current: translatePos,
      duration: duration / 1000, // Convert milliseconds to seconds
      ease: "expo.out",
      onUpdate: () => {
        const newStyle = `translateX(${elTranslatePos.current}px)`;
        gsap.set(carousel, {
          transform: newStyle,
          webkitTransform: newStyle
        });
      }
    });
  }
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div 
  bind:this={carousel}
  class="select-none relative"
  on:mousedown={startDragging}
  on:touchstart={startDragging}
>
  <div bind:this={contentDiv} class="flex cursor-grab active:cursor-grabbing">
    <slot {onLoad}></slot>
  </div>
</div>

<style global lang="postcss">
  @tailwind utilities;
  @tailwind components;
  @tailwind base;
</style>