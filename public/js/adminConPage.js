let ToggleState = true;
let InToggle = false;

function ShowPage() {
  $('#allContent').css('opacity', 1);
  $('#BG').css('opacity', 1);
}

function toggleOptionList() {
  let optionList = $('#optionList')[0];
  let toggleIcon = $('#optionList>#toggle>div>i')[0];
  if (InToggle) return;
  InToggle = true;
  if (ToggleState) {
    optionList.classList.add('hidden');
    toggleIcon.innerHTML = '&#xe84e;';
  } else {
    optionList.classList.remove('hidden');
    toggleIcon.innerHTML = '&#xe84f;';
  }
  setTimeout(() => {
    ToggleState = !ToggleState;
    InToggle = false;
    $(window).trigger('resize');
  }, 550);
}
