---
title: Surprising ease of gaming on Linux in 2020
date: 2020-06-07
---

During the quarantine my hobbies (and I imagine many others') shifted to more
indoors ones - possibly sedentary. Among other things I also started gaming
after a hiatus of many years.

But as you can probably imagine gaming on a 13" laptop with integrated graphics
is not the most gratifying of experiences. So I got myself a desktop with a
proper graphics card.

And now comes the interesting part: since all the games I originally planned on
playing have Linux support I decided to skip buying a copy of Windows and copy
my laptop [NixOS](https://nixos.org/) config and extend it as needed. Basically
I wanted to do some tinkering again.

# NVidia

The aforementioned graphics card is made by NVidia. And I braced for a world of
pain to set it up - the image of Linus flipping off the GPU giant still fresh in
my memory (actually it's been a while - that [happened in
2012](https://www.phoronix.com/scan.php?page=news_item&px=MTEyMTc))

![linus flipping off nvidia](/images/linus-nvidia.png)

And I recently heard a friend rant about NVidia driver / Steam compatibility
issues on Ubuntu. But the amazing Nix community had my back on this. Drives are
nicely packaged - I just needed to [enable
them](https://github.com/edofic/dotfiles/commit/b06aaa734bcc45e07e5c1e450494241aa11a4109#diff-c82c2d31640fd1352de397c99bfcda0cR191)

```nix
services.xserver.videoDrivers = [ "nvidia" ];
```

# Steam

Next up is my preferred distribution platform - Steam. Still have an account
(and some games) from my previous gaming sting. The steam installed also does
quite a bit of magic usually so this had me a bit worried as well. But
apparently I'm not the first first person (by far!) to try gaming on NixOS as
it's beautifully packaged. Just need to [add it to system
packages](https://github.com/edofic/dotfiles/commit/b06aaa734bcc45e07e5c1e450494241aa11a4109#diff-c82c2d31640fd1352de397c99bfcda0cR112)
and _everything_ works out of the box.

```nix
environment.systemPackages = with pkgs; [
  ...
  steam
  ...
];
```

Games install and play. I can get myself entertained, story over. Right?

# Steam Play / Proton / Wine

I incidentally learned about Steam Play (marketing name) / Proton (OSS name).
Valve created [a tool](https://github.com/ValveSoftware/Proton) to run
Windows-only games on Linux - basically Wine + (a lot of) tricks. And apparently
it works quite good. There is even [protondb.com](https://www.protondb.com/) - a
database of playability reports so you can tell what works and what doesn't
_before_ buying games.

And the best surprise? This is seamlessly integrated into Steam client. So I
picked I AAA (windows-only!) game I liked, clicked install, clicked play.

AND EVERYTHING WORKED!

I was running a AAA game via Wine on NixOS on a NVidia card. Butter smooth. I
was...still am...very impressed by this.

Even a random controller-only early access game worked through this setup. Very
nice. Speaking of controllers...

# PlayStation4 DualShock4

I tried some games on an old cheap controller that I had. It worked but it was
not great. I want to get myself something better. After a (literally) 2min search
I concluded that DS4 should work on Linux without too much trouble and I bought
one.

And it did work fine - over USB. Steam also comes with a driver for it comoplete
with configuration panel. Sorry about this post turning into a bit of an ad for
Steam but I am genuinely impressed.

But of course I wanted to go wireless. Proprietary drivers over bluetooth on
Linux on a funky distro. This will not end well. But I said I wanted tinkering.

I already had [bluetooth config](https://github.com/edofic/dotfiles/blob/afc0e990dd65c80b9233944d2af8ac19d637ee8e/nixos/desktop.nix#L55-L57) from my laptop

```nix
hardware.bluetooth = {
  enable = true;
  package = pkgs.bluezFull;
};
```

and a BT dongle for my PC. Almost a bit surprisingly it detected the controller
and successfully paired. But it showed up as a generic device. And here comes
the first bump on this journey: Steam drivers did not pick it up - and it was
unusable.

Long story short - after about an hour of research and trying things I added
[this udev rules](https://github.com/edofic/dotfiles/blob/afc0e990dd65c80b9233944d2af8ac19d637ee8e/nixos/desktop.nix#L55-L57) to enable Steam to have permissions to use the device.

```nix
services.udev.extraRules = ''
  # DualShock 4 over bluetooth hidraw
  KERNEL=="hidraw*", KERNELS=="*054C:05C4*", MODE="0666"
'';
```

And that's it. All the tinkering required. I can now use a DS4 over BT for AAA
Windows only games on NixOS. Based on my anecdotal evidence NixOS is now a
better (easier) plafform for gaming than Windows. But my criteria might be a bit
skewed ;)
